import Foundation

struct Article: Identifiable {
    let id = UUID()
    let title: String
    let summary: String
    let link: URL
    let source: String
    let author: String?
    let published: Date
    let term: [String]?
}

