<div class="bg-white rounded-xl shadow-sm p-6 mb-6 flex gap-6 items-start">
    <img src="{{.Country.Flag}}" alt="{{.Country.Name}} flag"
        class="w-32 h-24 object-cover rounded-lg border border-gray-100 flex-shrink-0" />
    <div class="flex-1">
        <span class="text-xs uppercase font-semibold text-rose-500 tracking-wider">{{.Country.Region}}</span>
        <h1 class="text-3xl font-bold text-gray-800 mt-1">{{.Country.Name}}</h1>
        <p class="text-gray-500 text-sm mb-4">{{.Country.Subregion}}</p>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
                <p class="text-gray-400 uppercase text-xs font-semibold">Capital</p>
                <p class="text-gray-700 font-medium">{{.Country.Capital}}</p>
            </div>
            <div>
                <p class="text-gray-400 uppercase text-xs font-semibold">Population</p>
                <p class="text-gray-700 font-medium">{{.Country.Population}}</p>
            </div>
            <div>
                <p class="text-gray-400 uppercase text-xs font-semibold">Currency</p>
                <p class="text-gray-700 font-medium">
                    {{range .Country.Currencies}}{{.}} {{end}}
                </p>
            </div>
            <div>
                <p class="text-gray-400 uppercase text-xs font-semibold">Languages</p>
                <p class="text-gray-700 font-medium">
                    {{range .Country.Languages}}{{.}} {{end}}
                </p>
            </div>
        </div>
    </div>
</div>

<!-- Add to Wishlist -->
{{if .IsLoggedIn}}
<div class="mb-6">
    <button id="add-wishlist-btn"
        onclick="addToWishlist('{{.Country.Name}}')"
        class="bg-rose-500 text-white px-5 py-2 rounded-lg text-sm font-semibold hover:bg-rose-600">
        + Add to Wishlist
    </button>
    <span id="wishlist-feedback" class="ml-3 text-sm text-gray-500"></span>
</div>
{{end}}

<!-- Attractions -->
<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div>
        <h2 class="text-xl font-bold text-gray-800 mb-3">Attractions &amp; Landmarks</h2>
        {{if .Attractions}}
            <ul class="space-y-2">
            {{range .Attractions}}
                <li class="bg-white rounded-lg px-4 py-3 shadow-sm text-sm">
                    <span class="font-medium text-gray-800">{{.Name}}</span>
                    <span class="text-gray-400 text-xs ml-2">{{.Kinds}}</span>
                </li>
            {{end}}
            </ul>
        {{else}}
            <p class="text-gray-400 text-sm">No attractions data available.</p>
        {{end}}
    </div>
</div>

<script src="/static/js/wishlist.js"></script>