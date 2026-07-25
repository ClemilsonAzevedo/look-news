import Foundation
import Combine

@MainActor
final class NewsStore: ObservableObject {
    static let shared = NewsStore()
    
    @Published var articles: [Article] = []
    @Published var isLoading = false
    @Published var errorMessage: String?
    @Published var lastUpdated: Date?
    
    private var pollingTask: Task<Void, Never>?
    
    private init() {
        startPolling()
    }
    
    func startPolling(interval: TimeInterval = 60 * 60 * 1) { // 1 hora
        pollingTask?.cancel()
        
        pollingTask = Task {
            await load()
            
            while !Task.isCancelled {
                try? await Task.sleep(nanoseconds: UInt64(interval * 1_000_000_000))
                await load()
            }
        }
    }
    
    func load() async {
        guard !isLoading else { return }
        
        isLoading = true
        defer { isLoading = false }
        
        do {
            let sources = SourcesStore.shared.sources
            let fetched = try await NewsService.shared.fetchArticles(sources: sources)
            articles = fetched
            errorMessage = nil
            lastUpdated = Date()
        } catch {
            errorMessage = "Não foi possível carregar os artigos."
            print("Erro ao buscar artigos: \(error)")
        }
    }
    
    func sourcesChanged() {
        Task {
            await load()
        }
    }
}
