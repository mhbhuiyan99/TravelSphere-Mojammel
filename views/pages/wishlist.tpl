<div class="mb-6">
    <h1 class="text-3xl font-bold text-gray-800">Travel Wishlist</h1>
    <p class="text-gray-500 mt-1 text-sm">Edit notes, update trip status, or remove destinations. Changes save without reloading the page.</p>
</div>

<div class="bg-white rounded-xl shadow-sm overflow-hidden">
    <table class="w-full text-sm">
        <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
                <th class="text-left px-4 py-3 text-xs font-semibold text-gray-400 uppercase">Country</th>
                <th class="text-left px-4 py-3 text-xs font-semibold text-gray-400 uppercase">Note</th>
                <th class="text-left px-4 py-3 text-xs font-semibold text-gray-400 uppercase">Status</th>
                <th class="text-left px-4 py-3 text-xs font-semibold text-gray-400 uppercase">Actions</th>
            </tr>
        </thead>
        <tbody id="wishlist-rows">
            {{range .Items}}
            <tr class="border-b border-gray-100" data-id="{{.ID}}">
                <td class="px-4 py-3 font-medium text-gray-700">{{.CountryName}}</td>
                <td class="px-4 py-3">
                    <input type="text" value="{{.Note}}" placeholder="Add a note..."
                        class="note-input border border-gray-200 rounded px-2 py-1 text-xs w-40 focus:outline-none focus:ring-1 focus:ring-rose-300" />
                </td>
                <td class="px-4 py-3">
                    <select class="status-select border border-gray-200 rounded px-2 py-1 text-xs focus:outline-none focus:ring-1 focus:ring-rose-300">
                        <option {{if eq .Status "Planned"}}selected{{end}}>Planned</option>
                        <option {{if eq .Status "Visited"}}selected{{end}}>Visited</option>
                    </select>
                </td>
                <td class="px-4 py-3 flex gap-2">
                    <button onclick="saveItem('{{.ID}}')"
                        class="bg-rose-500 text-white text-xs px-3 py-1 rounded hover:bg-rose-600">Save</button>
                    <button onclick="deleteItem('{{.ID}}')"
                        class="bg-gray-100 text-gray-600 text-xs px-3 py-1 rounded hover:bg-gray-200">Delete</button>
                </td>
            </tr>
            {{else}}
            <tr id="empty-row">
                <td colspan="4" class="px-4 py-6 text-center text-gray-400">No destinations yet. Browse countries and add some!</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</div>

<script src="/static/js/wishlist.js"></script>