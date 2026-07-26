import Foundation
import Combine
internal import SwiftUI

final class SourcesStore: ObservableObject {
    static let shared = SourcesStore()
    
    @Published var sources: [String] = []
    
    private let key = "looknews.sources"
    
    private init() {
        load()
    }
    
    func add(_ url: String) {
        let trimmed = url.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty,
              URL(string: trimmed) != nil,
              !sources.contains(trimmed) else { return }
        
        sources.append(trimmed)
        save()
        NewsStore.shared.sourcesChanged()
    }
    
    func remove(at offsets: IndexSet) {
        sources.remove(atOffsets: offsets)
        save()
        NewsStore.shared.sourcesChanged()
    }
    
    func remove(_ url: String) {
        sources.removeAll { $0 == url }
        save()
        NewsStore.shared.sourcesChanged()
    }
    
    private func save() {
        UserDefaults.standard.set(sources, forKey: key)
    }
    
    private func load() {
        sources = UserDefaults.standard.stringArray(forKey: key) ?? []
    }
}
