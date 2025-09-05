<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Recent treasures</title>
    <style>
        body { font-family: sans-serif; font-size: 12px; }
        h1 { text-align: center; margin-bottom: 20px; }
        .treasure { margin-bottom: 20px; border-bottom: 1px solid #ccc; padding-bottom: 10px; }
        .treasure:last-child { border-bottom: none; }
        img { margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Treasures Created</h1>

    @foreach($treasures as $treasure)
        <div class="treasure">
            <h2>{{ $treasure->name }}</h2>
            <p><strong>Point Gain:</strong> {{ $treasure->point_gain }}</p>
            <p><strong>Created at:</strong> {{ $treasure->created_at }}</p>

            @if($treasure->code?->image_url)
               <img src="{{ $treasure->code->image_url }}" width="120" alt="QR Code">
          @endif
               </div>
    @endforeach
</body>
</html>
