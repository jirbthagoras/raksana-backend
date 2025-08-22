<?php

namespace App\Filament\Resources\TreasureResource\Pages;

use App\Filament\Resources\TreasureResource;
use Filament\Actions;
use Filament\Resources\Pages\EditRecord;

class EditTreasure extends EditRecord
{
    protected static string $resource = TreasureResource::class;

    protected function getHeaderActions(): array
    {
        return [
            Actions\DeleteAction::make(),
        ];
    }
}
