<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Recent Quests</title>
    <style>
        body { font-family: sans-serif; font-size: 12px; }
        h1 { text-align: center; margin-bottom: 20px; }
        .quest { margin-bottom: 20px; border-bottom: 1px solid #ccc; padding-bottom: 10px; }
        .quest:last-child { border-bottom: none; }
        img { margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Quests Created in Last 3 Days</h1>

    @foreach($quests as $quest)
        <div class="quest">
            <h2>{{ $quest->detail->name }}</h2>
            <p><strong>Description:</strong> {{ $quest->detail->description }}</p>
            <p><strong>Point Gain:</strong> {{ $quest->detail->point_gain }}</p>
            <p><strong>Location:</strong> {{ $quest->location }}</p>
            <p><strong>Latitude:</strong> {{ $quest->latitude }}</p>
            <p><strong>Longitude:</strong> {{ $quest->longitude }}</p>
            <p><strong>Max Contributors:</strong> {{ $quest->max_contributors }}</p>
            <p><strong>Created at:</strong> {{ $quest->detail->created_at }}</p>

            @if($quest->latitude && $quest->longitude)
                <p><strong>Map:</strong></p>
                <img 
                    src="https://maps.googleapis.com/maps/api/staticmap?center={{ $quest->latitude }},{{ $quest->longitude }}&zoom=15&size=600x300&markers=color:red|{{ $quest->latitude }},{{ $quest->longitude }}&key={{ config('services.google.maps_key') }}" 
                    width="300" 
                    alt="Map">
            @endif

            <br>

            @if($quest->code?->image_url)
               <img src="{{ $quest->code->image_url }}" width="120" alt="QR Code">
            @endif
            </div>
    @endforeach
</body>
</html>
