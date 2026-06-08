const searchInput = document.getElementById('search-input');
const regionSelect = document.getElementById('region-select');
const results = document.getElementById('country-results');

function debounce(fn, delay) {
    let timer;
    return (...args) => {
        clearTimeout(timer);
        timer = setTimeout(() => fn(...args), delay);
    };
}

function buildCardsHTML(countries) {
    if (!countries.length) {
        return '<p class="text-gray-400">No countries found.</p>';
    }
    return `<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        ${countries.map(c => `
        <a href="/countries/${c.Slug}" class="bg-white rounded-xl shadow-sm overflow-hidden hover:shadow-md transition-shadow">
            <img src="${c.Flag}" alt="${c.Name} flag" class="w-full h-36 object-cover" />
            <div class="p-4">
                <h3 class="font-semibold text-gray-800">${c.Name}</h3>
                <p class="text-xs text-gray-500 mt-1"><span class="font-medium">Capital:</span> ${c.Capital}</p>
                <p class="text-xs text-gray-500"><span class="font-medium">Population:</span> ${c.Population.toLocaleString()}</p>
            </div>
        </a>`).join('')}
    </div>`;
}

async function fetchCountries() {
    const search = searchInput.value.trim();
    const region = regionSelect.value;

    results.innerHTML = '<p class="text-gray-400">Loading...</p>';

    try {
        const params = new URLSearchParams();
        if (search) params.set('search', search);
        if (region) params.set('region', region);

        const res = await fetch('/api/countries?' + params.toString());
        const json = await res.json();

        if (json.success) {
            results.innerHTML = buildCardsHTML(json.data);
        } else {
            results.innerHTML = '<p class="text-red-400">Failed to load countries.</p>';
        }
    } catch {
        results.innerHTML = '<p class="text-red-400">Network error. Please try again.</p>';
    }
}

const debouncedFetch = debounce(fetchCountries, 300);

searchInput.addEventListener('input', debouncedFetch);
regionSelect.addEventListener('change', fetchCountries);