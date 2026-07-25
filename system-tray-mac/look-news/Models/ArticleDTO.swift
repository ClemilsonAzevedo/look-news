import Foundation

struct ArticleDTO: Codable {
    let title: String
    let summary: String
    let link: String
    let date: String
    let source: String
    let author: String
    let published: String
    let terms: [String]?
}

extension ArticleDTO {
    func toDomain() -> Article? {
        guard let url = URL(string: link) else { return nil }

        return Article(
            title: title,
            summary: summary,
            link: url,
            source: source,
            author: author.isEmpty ? nil : author,
            published: ArticleDTO.parseDate(date),
            term: terms
        )
    }

    private static func parseDate(_ raw: String) -> Date {
        let withFraction = ISO8601DateFormatter()
        withFraction.formatOptions = [.withInternetDateTime, .withFractionalSeconds]

        let withoutFraction = ISO8601DateFormatter()
        withoutFraction.formatOptions = [.withInternetDateTime]

        return withFraction.date(from: raw) ?? withoutFraction.date(from: raw) ?? Date()
    }
}
