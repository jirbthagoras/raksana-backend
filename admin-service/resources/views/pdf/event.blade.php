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

            @if($event->code?->image_url)
               <img src="{{ $event->code->image_url }}" width="120" alt="QR Code">
          @endif
               </div>
    @endforeach
</body>
</html>
