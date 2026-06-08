<div class="mb-6">
    <h1 class="text-3xl font-bold text-gray-800">Country Explorer</h1>
    <p class="text-gray-500 mt-1">Browse every destination on first load. Search and filter update only the results below — no full page reload.</p>
</div>

<!-- Search & Filter -->
<div class="bg-white rounded-xl shadow-sm p-4 mb-6 flex gap-4">
    <div>
        <label class="text-xs text-gray-400 uppercase font-semibold block mb-1">Search</label>
        <input id="search-input" type="text" placeholder="Country or capital..."
            class="border border-gray-200 rounded px-3 py-2 text-sm w-56 focus:outline-none focus:ring-2 focus:ring-rose-300" />
    </div>
    <div>
        <label class="text-xs text-gray-400 uppercase font-semibold block mb-1">Region</label>
        <select id="region-select" class="border border-gray-200 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-rose-300">
            <option value="">All regions</option>
            <option>Africa</option>
            <option>Americas</option>
            <option>Asia</option>
            <option>Europe</option>
            <option>Oceania</option>
            <option>Antarctic</option>
        </select>
    </div>
</div>

<!-- Country Results — AJAX target -->
<div id="country-results">
    {{template "partials/country_cards.tpl" .}}
</div>

<script src="/static/js/countries.js"></script>