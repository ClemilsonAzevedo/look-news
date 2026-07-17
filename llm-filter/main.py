import json
import os
import sys
import time
from datetime import datetime, timezone
from email.utils import parsedate_to_datetime

from mlx_lm import load, stream_generate

FILTER_SYSTEM_PROMPT = """You are an extremely strict and precise article relevance filter.

Today's date: {today_date}
Criterion: {criterion}

You are given a numbered list of articles, each containing: title, summary, link, published_date, and terms.

**Core Rules (do not break any):**
- Accept an article ONLY if its MAIN subject directly and clearly matches the criterion.
- Reject anything tangential, loosely related, or where the topic is only mentioned in passing.
- Immediately reject ALL promotional content: ads, sponsored posts, press releases, "partnered with", product reviews, "best of" lists, affiliate marketing, or marketing-heavy articles.
- Prefer high-quality journalistic sources over blogs or SEO content.
- When multiple articles cover the same story, keep ONLY the most recent one.
- Strongly prioritize articles published today or in the last 48 hours.
- When in doubt, ALWAYS reject. Precision is far more important than recall.

**Output Format:**
Respond with ONLY a valid JSON object. Nothing else.

{
  "relevant": ["https://exact-link1.com", "https://exact-link2.com"]
}

Or if none are relevant:
{
  "relevant": []
}
"""

CLEANUP_SYSTEM_PROMPT = """You are performing the final consolidation pass on articles that already passed the initial relevance filter.

Today's date: {today_date}
Criterion: {criterion}

You now see the full list of candidates together, allowing you to detect duplicates across batches.

**Rules:**
- Re-evaluate every article strictly against the criterion. Drop anything that is only marginally relevant.
- Eliminate all promotional, sponsored, or low-quality content that may have slipped through.
- If two or more articles cover the same underlying story or event, keep ONLY the single most recent and highest-quality one.
- Strongly favor articles published today or very recently.
- When in doubt, reject.

**Output Format:**
Respond with ONLY a valid JSON object. No explanations, no extra text.

{
  "relevant": ["https://exact-link1.com", "https://exact-link2.com"]
}

Or:
{
  "relevant": []
}
"""

MODEL = os.environ.get("FILTER_MODEL", "prism-ml/Bonsai-8B-mlx-1bit")
BATCH_SIZE = int(os.environ.get("FILTER_BATCH_SIZE", "15"))
CLEANUP_BATCH_SIZE = int(os.environ.get("FILTER_CLEANUP_BATCH_SIZE", "25"))
MAX_TOKENS = int(os.environ.get("FILTER_MAX_TOKENS", "512"))
CLEANUP_MAX_TOKENS = int(os.environ.get("FILTER_CLEANUP_MAX_TOKENS", "768"))
SUMMARY_TRUNCATE = int(os.environ.get("FILTER_SUMMARY_TRUNCATE", "300"))
MAX_AGE_DAYS = int(os.environ.get("FILTER_MAX_AGE_DAYS", "1"))


def log(msg):
    print(f"[filter] {msg}", file=sys.stderr, flush=True)

def today_str():
    return datetime.now(timezone.utc).strftime("%Y-%m-%d")


log("loading model...")
model, tokenizer = load(MODEL)
log("ready")

def _parse_published(article):
    raw = article.get("published")
    if not raw:
        return None
    try:
        return datetime.fromisoformat(raw.replace("Z", "+00:00"))
    except ValueError:
        pass
    try:
        return parsedate_to_datetime(raw)
    except (TypeError, ValueError):
        return None

def _is_within_allowed_window(article):
    dt = _parse_published(article)
    if dt is None:
        return False
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    article_date = dt.astimezone(timezone.utc).date()
    today = datetime.now(timezone.utc).date()
    age_days = (today - article_date).days
    return 0 <= age_days <= MAX_AGE_DAYS

def _sort_by_recency(articles):
    def key(article):
        dt = _parse_published(article)
        if dt is None:
            return 1, 0
        if dt.tzinfo is None:
            dt = dt.replace(tzinfo=timezone.utc)
        return 0, -dt.timestamp()
    return sorted(articles, key=key)

def _format_article_line(index, article):
    title = article.get("title", "")
    summary = (article.get("summary", "") or "")[:SUMMARY_TRUNCATE]
    link = article.get("link", "")
    terms = ", ".join(article.get("terms") or [])
    published = article.get("published", "")
    text = f"[{index}] {title} - {summary}"
    if published:
        text += f" (published: {published})"
    if terms:
        text += f" (terms: {terms})"
    text += f"\n    link: {link}"
    return text

def _run_pass(system_prompt, query, batch, max_tokens, label, offset=0):
    listing = "\n".join(_format_article_line(offset + i, a) for i, a in enumerate(batch))
    user_content = f"TODAY: {today_str()}\nCRITERION: {query}\n\nARTICLES:\n{listing}"

    messages = [
        {"role": "system", "content": system_prompt},
        {"role": "user", "content": user_content},
    ]
    prompt = tokenizer.apply_chat_template(messages, add_generation_prompt=True)

    chunks = []
    n_tokens = 0
    depth = 0
    seen_open = False
    t0 = time.time()
    try:
        for resp in stream_generate(model, tokenizer, prompt, max_tokens=max_tokens):
            text = resp.text
            chunks.append(text)
            n_tokens += 1
            opens = text.count("{")
            closes = text.count("}")
            if opens:
                seen_open = True
                depth += opens
            if closes:
                depth -= closes
            if seen_open and depth <= 0:
                break
    except Exception as exc:
        log(f"[{label}] generation error: {exc}")
        return []
    dt = time.time() - t0
    rate = n_tokens / dt if dt > 0 else 0
    log(f"[{label}] {len(batch)} articles | {n_tokens} tokens | {dt:.1f}s | {rate:.1f} tok/s")

    output = "".join(chunks)
    valid_links = {a.get("link") for a in batch if a.get("link")}

    start = output.find("{")
    end = output.rfind("}")
    relevant = []
    if start != -1 and end != -1 and end > start:
        try:
            obj = json.loads(output[start:end + 1])
            for link in obj.get("relevant", []):
                if link in valid_links and link not in relevant:
                    relevant.append(link)
        except (ValueError, TypeError):
            log(f"[{label}] invalid json in output")
    else:
        log(f"[{label}] no json found in output")

    return relevant

def _first_pass(query, articles):
    by_link = {}
    order = []
    for offset in range(0, len(articles), BATCH_SIZE):
        batch = articles[offset:offset + BATCH_SIZE]
        n_batch = offset // BATCH_SIZE + 1
        links = _run_pass(FILTER_SYSTEM_PROMPT, query, batch, MAX_TOKENS,
                           f"filter {n_batch}", offset=offset)
        for link in links:
            if link not in by_link:
                article = next((a for a in batch if a.get("link") == link), None)
                if article is not None:
                    by_link[link] = article
                    order.append(link)
    return [by_link[link] for link in order]


def _cleanup_pass(query, candidates):
    if len(candidates) <= 1:
        return candidates

    kept_by_link = {}
    for offset in range(0, len(candidates), CLEANUP_BATCH_SIZE):
        chunk = candidates[offset:offset + CLEANUP_BATCH_SIZE]
        n_chunk = offset // CLEANUP_BATCH_SIZE + 1
        links = _run_pass(CLEANUP_SYSTEM_PROMPT, query, chunk, CLEANUP_MAX_TOKENS,
                           f"cleanup {n_chunk}", offset=offset)
        for link in links:
            article = next((a for a in chunk if a.get("link") == link), None)
            if article is not None:
                kept_by_link[link] = article

    if len(candidates) > CLEANUP_BATCH_SIZE:
        log("warning: candidates exceeded CLEANUP_BATCH_SIZE")

    return list(kept_by_link.values())

def filter_articles(query, articles):
    recent_articles = [a for a in articles if _is_within_allowed_window(a)]
    excluded = len(articles) - len(recent_articles)
    log(f"query={query!r} articles={len(articles)} in_window={len(recent_articles)} dropped={excluded}")

    if not recent_articles:
        return []

    first_pass = _first_pass(query, recent_articles)
    final = _cleanup_pass(query, first_pass)
    final = _sort_by_recency(final)
    log(f"done: {len(final)} relevant")

    return [a["link"] for a in final]

log("waiting for requests on stdin...")
for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        request = json.loads(line)
    except json.JSONDecodeError as exc:
        log(f"invalid json on stdin: {exc}")
        print(json.dumps({"relevant": []}), flush=True)
        continue
    try:
        relevant = filter_articles(request.get("query", ""), request.get("articles", []))
    except Exception as exc:
        log(f"unexpected error: {exc}")
        relevant = []
    print(json.dumps({"relevant": relevant}), flush=True)
    log("response sent")