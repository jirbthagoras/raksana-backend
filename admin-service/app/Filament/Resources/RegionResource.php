<?php

namespace App\Filament\Resources;

use App\Filament\Resources\RegionResource\Pages;
use App\Filament\Resources\RegionResource\RelationManagers;
use App\Models\Region;
use App\Models\Regions;
use Cheesegrits\FilamentGoogleMaps\Fields\Map;
use Filament\Forms;
use Filament\Forms\Components\TextInput;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;

class RegionResource extends Resource
{
    protected static ?string $model = Regions::class;

    protected static ?string $navigationIcon = 'heroicon-o-rectangle-stack';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                Forms\Components\TextInput::make("name")
                ->required(),
                Forms\Components\TextInput::make("location")
                ->required(),
                Map::make('map_picker')
                        ->label('Region Lintangfest Location (Pick on Map)')
                        ->defaultLocation([-6.175392, 106.827153])
                        ->draggable()
                        ->clickable()
                        ->afterStateUpdated(function ($state, callable $set) {
                            if (is_array($state) && count($state) >= 2) {
                                $set('latitude', $state["lat"]);
                                $set('longitude', $state["lng"]);
                            }
                        }),

                    Forms\Components\TextInput::make('latitude')
                        ->hidden()
                        ->dehydrated()
                        ->required(),

                    Forms\Components\TextInput::make('longitude')
                        ->hidden()
                        ->dehydrated()
                        ->required(),
            ]);
            
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                Tables\Columns\TextColumn::make("name"),
                Tables\Columns\TextColumn::make("location"),
                Tables\Columns\TextColumn::make("tree_amount"),
            ])
            ->filters([
                //
            ])
            ->actions([
                Tables\Actions\EditAction::make(),
            ])
            ->bulkActions([
                Tables\Actions\BulkActionGroup::make([
                    Tables\Actions\DeleteBulkAction::make(),
                ]),
            ]);
    }

    public static function getRelations(): array
    {
        return [
            //
        ];
    }

    public static function getPages(): array
    {
        return [
            'index' => Pages\ListRegions::route('/'),
            'create' => Pages\CreateRegion::route('/create'),
            'edit' => Pages\EditRegion::route('/{record}/edit'),
        ];
    }
}
