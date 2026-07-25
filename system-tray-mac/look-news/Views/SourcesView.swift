
internal import SwiftUI

struct SourcesView: View {
    @ObservedObject private var store = SourcesStore.shared
    @State private var newURL = ""
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Fontes RSS")
                .font(.headline)
            
            // Campo para adicionar
            HStack {
                TextField("https://...", text: $newURL)
                    .textFieldStyle(.roundedBorder)
                    .onSubmit(addSource)
                
                Button("Adicionar") {
                    addSource()
                }
                .disabled(newURL.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty)
            }
            
            // Lista de fontes
            if store.sources.isEmpty {
                Text("Nenhuma fonte adicionada")
                    .foregroundStyle(.secondary)
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                List {
                    ForEach(store.sources, id: \.self) { url in
                        HStack {
                            Text(url)
                                .lineLimit(1)
                                .truncationMode(.middle)
                            
                            Spacer()
                            
                            Button {
                                store.remove(url)
                            } label: {
                                Image(systemName: "trash")
                                    .foregroundStyle(.red)
                            }
                            .buttonStyle(.plain)
                        }
                    }
                    .onDelete(perform: store.remove)
                }
                .listStyle(.plain)
            }
        }
        .padding()
        .frame(width: 380, height: 320)
    }
    
    private func addSource() {
        store.add(newURL)
        newURL = ""
    }
}
