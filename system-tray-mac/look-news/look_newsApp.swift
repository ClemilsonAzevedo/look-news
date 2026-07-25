internal import SwiftUI

@main
struct LookNews: App {
    var body: some Scene {
        MenuBarExtra("look-news", systemImage: "eye") {
            ContentView()
        }
        .menuBarExtraStyle(.window)
    }
}
