<div class="mb-6">
    <h1 class="text-3xl font-bold text-gray-800">Dashboard</h1>
    <p class="text-gray-500 mt-1 text-sm">Your travel summary.</p>
</div>

<!-- Stats — AJAX target -->
<div id="dashboard-stats" class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
    <div class="bg-white rounded-xl shadow-sm p-6 text-center">
        <p class="text-4xl font-bold text-rose-500">{{.Total}}</p>
        <p class="text-gray-500 text-sm mt-1">Total Destinations</p>
    </div>
    <div class="bg-white rounded-xl shadow-sm p-6 text-center">
        <p class="text-4xl font-bold text-blue-500">{{.Planned}}</p>
        <p class="text-gray-500 text-sm mt-1">Planned</p>
    </div>
    <div class="bg-white rounded-xl shadow-sm p-6 text-center">
        <p class="text-4xl font-bold text-green-500">{{.Visited}}</p>
        <p class="text-gray-500 text-sm mt-1">Visited</p>
    </div>
</div>

<!-- Destination list -->
<div class="bg-white rounded-xl shadow-sm p-6">
    <h2 class="text-lg font-semibold text-gray-700 mb-4">Saved Destinations</h2>
    {{if .Items}}
    <ul class="space-y-2">
        {{range .Items}}
        <li class="flex items-center justify-between text-sm border-b border-gray-50 pb-2">
            <span class="font-medium text-gray-700">{{.CountryName}}</span>
            <span class="{{if eq .Status "Visited"}}text-green-500{{else}}text-blue-400{{end}} text-xs font-semibold">
                {{.Status}}
            </span>
        </li>
        {{end}}
    </ul>
    {{else}}
        <p class="text-gray-400 text-sm">No destinations saved yet.</p>
    {{end}}
</div>