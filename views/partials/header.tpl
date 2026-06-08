<nav class="bg-white border-b border-gray-200 px-6 py-3 flex items-center gap-8">
    <a href="/" class="text-rose-500 font-bold text-xl">TravelSphere</a>

    <ul class="flex gap-6 list-none">
        <li><a href="/" class="text-sm {{if eq .CurrentPath "/"}}text-rose-500 font-semibold{{else}}text-gray-500 hover:text-gray-800{{end}}">Home</a></li>
        <li><a href="/countries" class="text-sm {{if eq .CurrentPath "/countries"}}text-rose-500 font-semibold{{else}}text-gray-500 hover:text-gray-800{{end}}">Countries</a></li>
        <li><a href="/wishlist" class="text-sm {{if eq .CurrentPath "/wishlist"}}text-rose-500 font-semibold{{else}}text-gray-500 hover:text-gray-800{{end}}">Wishlist</a></li>
        <li><a href="/dashboard" class="text-sm {{if eq .CurrentPath "/dashboard"}}text-rose-500 font-semibold{{else}}text-gray-500 hover:text-gray-800{{end}}">Dashboard</a></li>
    </ul>

    <div class="ml-auto flex items-center gap-4 text-sm">
        {{if .IsLoggedIn}}
            <span class="text-gray-600">Hi, <strong>{{.Username}}</strong></span>
            <a href="/logout" class="text-rose-500 hover:underline">Logout</a>
        {{else}}
            <a href="/login" class="bg-rose-500 text-white px-4 py-1.5 rounded hover:bg-rose-600">Login</a>
        {{end}}
    </div>
</nav>