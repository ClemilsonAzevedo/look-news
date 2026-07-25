internal import SwiftUI

struct ArticleRow: View {
    let article: Article
    @State private var isHovering = false

    private let accent = Color.brandAccent

    var body: some View {
        Button {
            NSWorkspace.shared.open(article.link)
        } label: {
            HStack(alignment: .top, spacing: 10) {
                VStack(alignment: .leading, spacing: 5) {
                    HStack(spacing: 6) {
                        Text(article.source.uppercased())
                            .font(.system(size: 10, weight: .bold))
                            .tracking(0.4)
                            .foregroundStyle(accent)
                            .lineLimit(1)

                        Circle()
                            .fill(.tertiary)
                            .frame(width: 2.5, height: 2.5)

                        Text(article.published, style: .relative)
                            .font(.system(size: 10.5))
                            .foregroundStyle(.tertiary)
                    }

                    Text(article.title)
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundStyle(.primary)
                        .lineLimit(2)
                        .multilineTextAlignment(.leading)
                        .fixedSize(horizontal: false, vertical: true)

                    Text(article.summary)
                        .font(.system(size: 11.5))
                        .foregroundStyle(.secondary)
                        .lineLimit(1)
                        .multilineTextAlignment(.leading)
                }

                Spacer(minLength: 0)

                Image(systemName: "arrow.up.right")
                    .font(.system(size: 10, weight: .semibold))
                    .foregroundStyle(accent)
                    .opacity(isHovering ? 1 : 0)
                    .offset(x: isHovering ? 0 : -4)
                    .padding(.top, 2)
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 11)
            .background {
                RoundedRectangle(cornerRadius: 10, style: .continuous)
                    .fill(isHovering ? accent.opacity(0.10) : .clear)
                    .glassEffect(isHovering ? .regular.tint(accent.opacity(0.14)).interactive() : .identity, in: RoundedRectangle(cornerRadius: 10, style: .continuous))
            }
            .padding(.horizontal, 6)
            .contentShape(Rectangle())
        }
        .buttonStyle(.plain)
        .animation(.easeOut(duration: 0.15), value: isHovering)
        .onHover { hovering in
            isHovering = hovering
        }
    }
}
