self.addEventListener("install", (event) => {
    console.log("service worker installed");
});

self.addEventListener("activate", (event) => {
    console.log("Service worker activated");
    // This ensures the service worker activates right away
    event.waitUntil(self.clients.claim());
});


self.addEventListener("push", (event) => {
    const data = event.data.json();
    console.log(data);
    const options = {
        body: data.body,
        icon: data.icon,
        tag: data.tag,
        silent: data.silent ?? true
    }
    event.waitUntil(
        self.registration.showNotification(data.title, options)
    );
});
