// Add to wishlist from destination page
async function addToWishlist(countryName) {
    const feedback = document.getElementById('wishlist-feedback');
    const btn = document.getElementById('add-wishlist-btn');
    if (btn) btn.disabled = true;

    try {
        const res = await fetch('/api/wishlist', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ country_name: countryName }),
        });
        const json = await res.json();
        if (json.success) {
            feedback.textContent = '✓ Added to wishlist';
            feedback.className = 'ml-3 text-sm text-green-500';
        } else {
            feedback.textContent = json.message || 'Could not add.';
            feedback.className = 'ml-3 text-sm text-red-400';
            if (btn) btn.disabled = false;
        }
    } catch {
        feedback.textContent = 'Network error.';
        feedback.className = 'ml-3 text-sm text-red-400';
        if (btn) btn.disabled = false;
    }
}

// Save note + status update from wishlist page
async function saveItem(id) {
    const row = document.querySelector(`tr[data-id="${id}"]`);
    const note = row.querySelector('.note-input').value;
    const status = row.querySelector('.status-select').value;

    const res = await fetch(`/api/wishlist/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ note, status }),
    });
    const json = await res.json();
    if (json.success) {
        refreshDashboardStats();
    } else {
        alert(json.message || 'Update failed.');
    }
}

// Delete item from wishlist page
async function deleteItem(id) {
    const res = await fetch(`/api/wishlist/${id}`, { method: 'DELETE' });
    const json = await res.json();
    if (json.success) {
        document.querySelector(`tr[data-id="${id}"]`).remove();
        refreshDashboardStats();
    } else {
        alert(json.message || 'Delete failed.');
    }
}

// Refresh dashboard stat counters without page reload
async function refreshDashboardStats() {
    const statsDiv = document.getElementById('dashboard-stats');
    if (!statsDiv) return;

    const res = await fetch('/api/dashboard/summary');
    const json = await res.json();
    if (!json.success) return;

    const { total, planned, visited } = json.data;
    statsDiv.innerHTML = `
        <div class="bg-white rounded-xl shadow-sm p-6 text-center">
            <p class="text-4xl font-bold text-rose-500">${total}</p>
            <p class="text-gray-500 text-sm mt-1">Total Destinations</p>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 text-center">
            <p class="text-4xl font-bold text-blue-500">${planned}</p>
            <p class="text-gray-500 text-sm mt-1">Planned</p>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 text-center">
            <p class="text-4xl font-bold text-green-500">${visited}</p>
            <p class="text-gray-500 text-sm mt-1">Visited</p>
        </div>`;
}