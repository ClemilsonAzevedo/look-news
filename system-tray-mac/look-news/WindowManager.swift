internal import SwiftUI
import AppKit

final class WindowManager: NSObject {
    static let shared = WindowManager()
    
    private var aboutWindow: NSWindow?
    private var sourcesWindow: NSWindow?
    
    func showAbout() {
        if aboutWindow == nil {
            let window = makeWindow(title: "Sobre o look news", size: NSSize(width: 320, height: 280))
            window.contentView = NSHostingView(rootView: AboutView())
            aboutWindow = window
        }
        present(aboutWindow)
    }
    
    func showSources() {
        if sourcesWindow == nil {
            let window = makeWindow(title: "Gerenciar Fontes", size: NSSize(width: 400, height: 360))
            window.contentView = NSHostingView(rootView: SourcesView())
            sourcesWindow = window
        }
        present(sourcesWindow)
    }
    
    private func makeWindow(title: String, size: NSSize) -> NSWindow {
        let window = NSWindow(
            contentRect: NSRect(origin: .zero, size: size),
            styleMask: [.titled, .closable],
            backing: .buffered,
            defer: false
        )
        window.title = title
        window.isReleasedWhenClosed = false
        window.center()
        return window
    }
    
    private func present(_ window: NSWindow?) {
        NSApp.activate(ignoringOtherApps: true)
        window?.makeKeyAndOrderFront(nil)
    }
}
