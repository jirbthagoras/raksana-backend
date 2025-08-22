<?php

namespace App\Filament\Resources\TreasureResource\Pages;

use App\Filament\Resources\TreasureResource;
use Filament\Actions;
use Filament\Resources\Pages\ListRecords;

class ListTreasures extends ListRecords
{
    protected static string $resource = TreasureResource::class;

    protected function getHeaderActions(): array
    {
        return [
            Actions\CreateAction::make(),
        ];
    }
}
