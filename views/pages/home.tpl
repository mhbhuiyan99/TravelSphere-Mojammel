<!-- Hero -->
<div class="rounded-2xl bg-gradient-to-br from-gray-800 to-gray-900 text-white px-10 py-16 mb-10">
    <h1 class="text-4xl font-bold mb-2">Discover your next destination</h1>
    <p class="text-gray-300 mb-6">Search countries, explore attractions, and curate your personal travel wishlist.</p>

    <!-- Search with autocomplete -->
    <div class="relative w-full max-w-md">
        <label class="text-xs uppercase text-gray-400 font-semibold block mb-1">Where to next?</label>
        <input id="home-search" type="text" placeholder="Country or capital..."
            autocomplete="off"
            class="w-full px-4 py-2 rounded-lg text-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-rose-400" />
        <div id="autocomplete-list"
            class="absolute z-10 bg-white border border-gray-200 rounded-lg shadow-lg w-full mt-1 hidden text-sm text-gray-700">
        </div>
    </div>
</div>

<!-- Featured Destinations -->
<h2 class="text-xl font-bold text-gray-800 mb-4">Featured destinations</h2>
<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-10">
    {{range .Featured}}
    <a href="/countries/{{.Slug}}"
        class="bg-white rounded-xl shadow-sm overflow-hidden hover:shadow-md transition-shadow">
        <img src="{{.Flag}}" alt="{{.Name}}" class="w-full h-28 object-cover" />
        <div class="p-3">
            <p class="font-semibold text-gray-800 text-sm">{{.Name}}</p>
            <p class="text-xs text-gray-400">{{.Capital}} · {{.Region}}</p>
        </div>
    </a>
    {{end}}
</div>

<!-- Popular Attractions -->
<h2 class="text-xl font-bold text-gray-800 mb-4">Popular attractions</h2>
<div class="bg-white rounded-xl shadow-sm divide-y divide-gray-100">
    {{range .Attractions}}
    <div class="px-5 py-3 flex items-center justify-between">
        <span class="text-sm font-medium text-gray-700">{{.Name}}</span>
        <span class="text-xs text-gray-400">{{.Kinds}}</span>
    </div>
    {{else}}
    <p class="px-5 py-4 text-gray-400 text-sm">Attractions data unavailable.</p>
    {{end}}
</div>

<script src="/static/js/autocomplete.js"></script>