<?php

namespace App\Filament\Resources\RegionResource\Pages;

use App\Filament\Resources\RegionResource;
use Filament\Actions;
use Filament\Resources\Pages\CreateRecord;

class CreateRegion extends CreateRecord
{
    protected static string $resource = RegionResource::class;

    protected function mutateFormDataBeforeCreate(array $data): array
    {
        if (isset($data['map_picker']['lat'], $data['map_picker']['lng'])) {
            $data['latitude'] = $data['map_picker']['lat'];
            $data['longitude'] = $data['map_picker']['lng'];
        }

        unset($data['map_picker']);

        return $data;
    }
}
