<?php

namespace App\Filament\Resources\QuestResource\Pages;

use App\Filament\Resources\QuestResource;
use App\Models\Codes;
use App\Models\Details;
use Endroid\QrCode\Builder\Builder;
use Endroid\QrCode\Encoding\Encoding;
use Endroid\QrCode\Color\Color;
use Filament\Resources\Pages\CreateRecord;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Str;
use Firebase\JWT\JWT;

class CreateQuest extends CreateRecord
{
    protected static string $resource = QuestResource::class;

    protected function mutateFormDataBeforeCreate(array $data): array
    {
        $detail = Details::create([
            'name' => $data['detail_name'],
            'description' => $data['detail_description'],
            'point_gain' => $data['detail_point_gain'],
        ]);

        $data['detail_id'] = $detail->id;

        unset($data['detail_name'], $data['detail_description'], $data['detail_point_gain']);

        $uuid = (string) Str::uuid();

        $payload = [
            'uuid' => $uuid,
            'exp' => time() + 3600,
        ];

        $secretKey = env('JWT_SECRET_KEY', env("JWT_SECRET_KEY"));
        $jwt = JWT::encode($payload, $secretKey, 'HS256');

        $result = Builder::create()
            ->data($jwt)
            ->encoding(new Encoding('UTF-8'))
            ->size(50)
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
        ]);

        $data['code_id'] = $code->id;

        return $data;
    }
}
