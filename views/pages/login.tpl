<div class="max-w-md mx-auto mt-16 bg-white rounded-xl shadow p-8">
    <h1 class="text-2xl font-bold text-gray-800 mb-2">Welcome to TravelSphere</h1>
    <p class="text-gray-500 text-sm mb-6">Enter your name to start planning your trips.</p>

    {{if .Error}}
        <div class="bg-red-50 text-red-600 text-sm px-4 py-2 rounded mb-4">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/login" class="flex flex-col gap-4">
        <input
            type="text"
            name="username"
            placeholder="Your name"
            autofocus
            class="border border-gray-300 rounded px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-rose-400"
        />
        <button
            type="submit"
            class="bg-rose-500 text-white rounded px-4 py-2 text-sm font-semibold hover:bg-rose-600"
        >
            Enter
        </button>
    </form>
</div>