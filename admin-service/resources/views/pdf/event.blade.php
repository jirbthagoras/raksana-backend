<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Recent events</title>
    <style>
        body { font-family: sans-serif; font-size: 12px; }
        h1 { text-align: center; margin-bottom: 20px; }
        .event { margin-bottom: 20px; border-bottom: 1px solid #ccc; padding-bottom: 10px; }
        .event:last-child { border-bottom: none; }
        img { margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Events Created in Last Week</h1>

    @foreach($events as $event)
        <div class="event">
            <h2>{{ $event->detail->name }}</h2>
            <p><strong>Description:</strong> {{ $event->detail->description }}</p>
            <p><strong>Point Gain:</strong> {{ $event->detail->point_gain }}</p>
            <p><strong>Location:</strong> {{ $event->location }}</p>
            <p><strong>Latitude:</strong> {{ $event->latitude }}</p>
            <p><strong>Longitude:</strong> {{ $event->longitude }}</p>
            <p><strong>Contact Person:</strong> {{ $event->contact }}</p>
            <p><strong>Starts at:</strong> {{ $event->starts_at }}</p>
            <p><strong>Ends at:</strong> {{ $event->ends_at }}</p>
            <p><strong>Created at:</strong> {{ $event->details->created_at }}</p>

            @if($event->latitude && $event->longitude)
                <p><strong>Map:</strong></p>
                <img 
                    src="https://maps.googleapis.com/maps/api/staticmap?center={{ $event->latitude }},{{ $event->longitude }}&zoom=15&size=600x300&markers=color:red|{{ $event->latitude }},{{ $event->longitude }}&key={{ config('services.google.maps_key') }}" 
                    width="300" 
                    alt="Map">
            @endif
            <br>

            @if($event->code?->image_url)
               <img src="{{ $event->code->image_url }}" width="120" alt="QR Code">
            @endif
        </div>
    @endforeach
</body>
</html>
