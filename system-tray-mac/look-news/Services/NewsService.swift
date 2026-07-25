import Foundation

enum APIError: Error {
    case invalidURL
    case invalidResponse
    case decodingFailed
}

final class NewsService {
    static let shared = NewsService()
    
    private let baseURL = "http://127.0.0.1:8080"
    
    func fetchArticles(sources: [String]) async throws -> [Article] {
        guard let url = URL(string: "\(baseURL)/news") else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        request.httpBody = try JSONEncoder().encode(sources)
        
        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw APIError.invalidResponse
        }
        
        do {
            let dtos = try JSONDecoder().decode([ArticleDTO].self, from: data)
            return dtos.compactMap { $0.toDomain() }
        } catch {
            print("Erro de decoding: \(error)")
            throw APIError.decodingFailed
        }
    }
}
