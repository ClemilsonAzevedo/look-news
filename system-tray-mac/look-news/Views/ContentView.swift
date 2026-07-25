internal import SwiftUI

struct ContentView: View {
    @ObservedObject private var store = NewsStore.shared
    
    var body: some View {
        VStack(spacing: 0) {
            header
            Divider()
            articleList
            Divider()
            footer
        }
        .frame(width: 380, height: 520)
        .background(.regularMaterial)
    }
    
    // MARK: - Header
    
    private var header: some View {
        HStack(spacing: 10) {
            HStack(spacing: 6) {
                Image(systemName: "eye.fill")
                    .font(.system(size: 13, weight: .semibold))
                Text("look news")
                    .font(.system(size: 13, weight: .semibold))
            }
            .foregroundStyle(.primary)
            
            Spacer()
            
            GlassEffectContainer(spacing: 8) {
                HStack(spacing: 8) {
                    Button {
                        refresh()
                    } label: {
                        Image(systemName: "arrow.clockwise")
                            .font(.system(size: 12, weight: .medium))
                            .rotationEffect(.degrees(store.isLoading ? 360 : 0))
                            .frame(width: 28, height: 28)
                    }
                    .buttonStyle(.plain)
                    .foregroundStyle(.secondary)
                    .glassEffect(.regular.interactive(), in: .circle)
                    .animation(
                        store.isLoading
                            ? .linear(duration: 0.7).repeatForever(autoreverses: false)
                            : .default,
                        value: store.isLoading
                    )
                    .disabled(store.isLoading)
                    
                    Button {
                        WindowManager.shared.showSources()
                    } label: {
                        Image(systemName: "link")
                            .font(.system(size: 12, weight: .medium))
                            .frame(width: 28, height: 28)
                    }
                    .buttonStyle(.plain)
                    .foregroundStyle(.secondary)
                    .glassEffect(.regular.interactive(), in: .circle)
                    
                    Button {
                        NSApplication.shared.terminate(nil)
                    } label: {
                        Image(systemName: "power")
                            .font(.system(size: 12, weight: .medium))
                            .frame(width: 28, height: 28)
                    }
                    .buttonStyle(.plain)
                    .foregroundStyle(.red)
                    .glassEffect(.regular.tint(.red.opacity(0.15)).interactive(), in: .circle)
                }
            }
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 10)
    }
    
    // MARK: - Lista
    
    private var articleList: some View {
        ScrollView {
            LazyVStack(spacing: 2) {
                if store.isLoading && store.articles.isEmpty {
                    ProgressView()
                        .controlSize(.small)
                        .padding(.vertical, 32)
                } else if let errorMessage = store.errorMessage {
                    VStack(spacing: 6) {
                        Image(systemName: "wifi.exclamationmark")
                            .font(.system(size: 16))
                            .foregroundStyle(.secondary)
                        Text(errorMessage)
                            .font(.system(size: 11))
                            .foregroundStyle(.secondary)
                            .multilineTextAlignment(.center)
                    }
                    .padding(.vertical, 32)
                } else if store.articles.isEmpty {
                    Text("Nenhum artigo por aqui ainda.")
                        .font(.system(size: 11))
                        .foregroundStyle(.secondary)
                        .padding(.vertical, 32)
                } else {
                    ForEach(store.articles) { article in
                        ArticleRow(article: article)
                    }
                }
            }
            .padding(.vertical, 4)
        }
        .scrollIndicators(.hidden)
    }
    
    // MARK: - Footer
    
    private var footer: some View {
        HStack(spacing: 6) {
            Text("\(store.articles.count) artigos")
                .font(.system(size: 11, weight: .medium))
                .foregroundStyle(.secondary)
            
            Circle()
                .fill(.tertiary)
                .frame(width: 2.5, height: 2.5)
            
            Text(lastUpdatedText)
                .font(.system(size: 11))
                .foregroundStyle(.tertiary)
            
            Spacer()
            
            Button {
                WindowManager.shared.showAbout()
            } label: {
                Image(systemName: "info.circle")
                    .font(.system(size: 13, weight: .medium))
                    .frame(width: 24, height: 24)
            }
            .buttonStyle(.plain)
            .foregroundStyle(.secondary)
            .glassEffect(.regular.interactive(), in: .circle)
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 10)
    }
    
    private var lastUpdatedText: String {
        guard let date = store.lastUpdated else {
            return "atualizando..."
        }
        
        let formatter = RelativeDateTimeFormatter()
        formatter.unitsStyle = .short
        return "atualizado \(formatter.localizedString(for: date, relativeTo: Date()))"
    }
    
    // MARK: - Ações
    
    private func refresh() {
        Task {
            await store.load()
        }
    }
}
