<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} — TravelSphere</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50 text-gray-800 min-h-screen flex flex-col">
    {{template "partials/header.tpl" .}}

    <main class="flex-1 max-w-7xl mx-auto w-full px-4 py-8">
        {{.LayoutContent}}
    </main>

    {{template "partials/footer.tpl" .}}
</body>
</html>