<?php

namespace App\Filament\Resources\TreasureResource\Pages;

use App\Filament\Resources\TreasureResource;
use App\Models\Codes;
use Endroid\QrCode\Builder\Builder;
use Endroid\QrCode\Encoding\Encoding;
use Endroid\QrCode\Color\Color;
use Filament\Resources\Pages\CreateRecord;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Str;
use Firebase\JWT\JWT;

class CreateTreasure extends CreateRecord
{
    protected static string $resource = TreasureResource::class;

    protected function mutateFormDataBeforeCreate(array $data): array
    {
        $uuid = (string) Str::random(12);

        $payload = [
            'uuid' => $uuid,
            'type' => "treasure",
            'exp'  => time() + (365 * 24 * 60 * 60),
        ];


        $secretKey = env('JWT_SECRET_KEY', env("JWT_SECRET_KEY"));
        $jwt = JWT::encode($payload, $secretKey, 'HS256');

        $result = Builder::create()
            ->data($jwt)
            ->encoding(new Encoding('UTF-8'))
            ->size(100)
            ->margin(5)
            ->foregroundColor(new Color(0, 0, 0))
            ->backgroundColor(new Color(255, 255, 255))
            ->build();

        $qrImage = $result->getString();

        $path = "qr/{$uuid}.png";
        Storage::disk('s3')->put($path, $qrImage);
        $url = Storage::disk('s3')->url($path);

        $code = Codes::create([
            'id' => $uuid,
            'image_url' => $url,
            'jwt' => $jwt,
        ]);

        $data['code_id'] = $code->id;

        return $data;
    }
}
