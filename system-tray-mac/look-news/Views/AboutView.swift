internal import SwiftUI

struct AboutView: View {
    var body: some View {
        VStack(spacing: 12) {
            Image(systemName: "eye.fill")
                .font(.system(size: 36))
                .foregroundStyle(.blue)
                .padding(.top, 20)

            Text("look news")
                .font(.system(size: 16, weight: .semibold))

            Text("Versão 1.0.0")
                .font(.system(size: 12))
                .foregroundStyle(.secondary)

            Divider()
                .padding(.horizontal, 40)

            Text("Agregador de notícias direto na sua barra de menu.")
                .font(.system(size: 12))
                .foregroundStyle(.secondary)
                .multilineTextAlignment(.center)
                .padding(.horizontal, 30)

            Spacer()

            Text("© 2026 Clemilson de Azevedo")
                .font(.system(size: 10))
                .foregroundStyle(.tertiary)
                .padding(.bottom, 16)
        }
        .frame(width: 320, height: 280)
    }
}

//struct UpdateCheckView: View {
//    @State private var isChecking = true
//    @State private var isUpToDate = false
//
//    var body: some View {
//        VStack(spacing: 16) {
//            if isChecking {
//                ProgressView()
//                    .controlSize(.small)
//                Text("Verificando atualizações...")
//                    .font(.system(size: 12))
//                    .foregroundStyle(.secondary)
//            } else {
//                Image(systemName: "checkmark.circle.fill")
//                    .font(.system(size: 32))
//                    .foregroundStyle(.green)
//                Text("Você está usando a versão mais recente")
//                    .font(.system(size: 12, weight: .medium))
//                Text("look news 1.0.0")
//                    .font(.system(size: 11))
//                    .foregroundStyle(.secondary)
//            }
//        }
//        .frame(width: 320, height: 200)
//        .task {
//            try? await Task.sleep(nanoseconds: 900_000_000)
//            isChecking = false
//            isUpToDate = true
//        }
//    }
//}
