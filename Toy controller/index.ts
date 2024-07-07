const tokenServer = new URL('http://[::1]:3001');

Bun.serve({
    hostname: '::',
    port: 3000,
    async fetch(req) {
        const url = new URL(req.url)
        console.debug(`[${new Date().toLocaleString()}]:`, url.pathname + url.search)

        switch (true) {
            case url.pathname === '/login' && url.search === '':
                return new Response(Bun.file('login.html'));

            case url.pathname === '/login' && url.searchParams.has('username') && url.searchParams.has('password'):
                const qrPageContent = await generateQRPage(url.searchParams.get('username')!, url.searchParams.get('password')!);
                return new Response(new File([qrPageContent], 'qr.html', { type: "text/html" }));

            case url.pathname === '/favicon.ico':
                return new Response(Bun.file('image.png'))

            default:
                const headers = new Headers();
                headers.set('Location', '/login')
                return new Response(null, { headers, status: 303 });
        }
    }
});

async function generateQRPage(username: string, password: string) {
    const params = new URLSearchParams();
    params.set('username', username);
    params.set('password', password);
    const url = new URL(`/qr?${params.toString()}`, tokenServer)

    const qrPageTemplate = await Bun.file('qr.html').text();
    const qrLink = await (await fetch(url)).text();

    return new HTMLRewriter()
        .on('#qr', {
            element(element) {
                element.setAttribute('src', qrLink);
            },
        })
        .transform(qrPageTemplate);
}
