<?php

namespace App\Filament\Resources;

use App\Filament\Resources\TreasureResource\Pages;
use App\Models\Treasure;
use App\Models\Treasures;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Tables;
use Filament\Tables\Table;
use Filament\Resources\Resource;

class TreasureResource extends Resource
{
    protected static ?string $model = Treasures::class;

    protected static ?string $navigationIcon = 'heroicon-o-gift';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                Forms\Components\TextInput::make('name')->required(),
                Forms\Components\TextInput::make('point_gain')->numeric()->required(),
                Forms\Components\Toggle::make('claimed')->default(false),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                Tables\Columns\TextColumn::make('name'),
                Tables\Columns\TextColumn::make('point_gain'),
                Tables\Columns\IconColumn::make('claimed')->boolean(),
                Tables\Columns\ImageColumn::make('code.image_url')->label('QR Code'),
            ])
            ->actions([
                Tables\Actions\EditAction::make(),
            ])
            ->bulkActions([
                Tables\Actions\DeleteBulkAction::make(),
            ]);
    }

    public static function getPages(): array
    {
        return [
            'index' => Pages\ListTreasures::route('/'),
            'create' => Pages\CreateTreasure::route('/create'),
            'edit' => Pages\EditTreasure::route('/{record}/edit'),
        ];
    }
}
