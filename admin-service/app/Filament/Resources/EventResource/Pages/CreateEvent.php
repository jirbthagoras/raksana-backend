<?php

namespace App\Filament\Resources\EventResource\Pages;

use App\Filament\Resources\EventResource;
use App\Models\Codes;
use App\Models\Details;
use Endroid\QrCode\Builder\Builder;
use Endroid\QrCode\Encoding\Encoding;
use Endroid\QrCode\Color\Color;
use Filament\Resources\Pages\CreateRecord;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Str;
use Firebase\JWT\JWT;

class CreateEvent extends CreateRecord
{
    protected static string $resource = EventResource::class;

    protected function mutateFormDataBeforeCreate(array $data): array
    {
        $detail = Details::create([
            'name' => $data['detail_name'],
            'description' => $data['detail_description'],
            'point_gain' => $data['detail_point_gain'],
        ]);
        $data['detail_id'] = $detail->id;

        unset($data['detail_name'], $data['detail_description'], $data['detail_point_gain']);

        if (isset($data['map_picker']['lat'], $data['map_picker']['lng'])) {
            $data['latitude'] = $data['map_picker']['lat'];
            $data['longitude'] = $data['map_picker']['lng'];
        }

        unset($data['map_picker']);

        $data['detail_id'] = $detail->id;

        unset($data['detail_name'], $data['detail_description'], $data['detail_point_gain']);

        $uuid = (string) Str::uuid();

        $nbf = strtotime($data['starts_at']);
        $exp = strtotime($data['ends_at']);

        if ($nbf === $exp) {
            $exp = $nbf + (3 * 24 * 60 * 60);
        }

        $payload = [
            'uuid' => $uuid,
            'nbf'  => $nbf,
            'exp'  => $exp,
            'type' => "event",
        ];

        $secretKey = env('JWT_SECRET_KEY', env("JWT_SECRET_KEY"));
        $jwt = JWT::encode($payload, $secretKey, 'HS256');

        $result = Builder::create()
            ->data($jwt)
            ->encoding(encoding: new Encoding('UTF-8'))
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

        // cover_path already handled by FileUpload
        return $data;
    }
}
