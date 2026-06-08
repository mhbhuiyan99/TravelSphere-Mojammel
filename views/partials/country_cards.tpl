{{if .Error}}
    <p class="text-red-500">{{.Error}}</p>
{{else if not .Countries}}
    <p class="text-gray-400">No countries found.</p>
{{else}}
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {{range .Countries}}
        <a href="/countries/{{.Slug}}" class="bg-white rounded-xl shadow-sm overflow-hidden hover:shadow-md transition-shadow">
            <img src="{{.Flag}}" alt="{{.Name}} flag" class="w-full h-36 object-cover" />
            <div class="p-4">
                <h3 class="font-semibold text-gray-800">{{.Name}}</h3>
                <p class="text-xs text-gray-500 mt-1"><span class="font-medium">Capital:</span> {{.Capital}}</p>
                <p class="text-xs text-gray-500"><span class="font-medium">Population:</span> {{.Population}}</p>
                <p class="text-xs text-gray-500"><span class="font-medium">Currencies:</span> {{join .Currencies ", "}}</p>
                <p class="text-xs text-gray-500"><span class="font-medium">Languages:</span> {{join .Languages ", "}}</p>
            </div>
        </a>
        {{end}}
    </div>
{{end}}