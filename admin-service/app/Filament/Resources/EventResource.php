<?php

namespace App\Filament\Resources;

use App\Filament\Resources\EventResource\Pages;
use App\Models\Event;
use Filament\Forms;
use Filament\Tables;
use Filament\Resources\Resource;
use Illuminate\Support\Facades\Storage;

class EventResource extends Resource
{
    protected static ?string $model = Event::class;
    protected static ?string $navigationIcon = 'heroicon-o-calendar';

    public static function form(Forms\Form $form): Forms\Form
    {
        return $form
            ->schema([
                Forms\Components\Section::make('Event Detail')
                    ->schema([
                        Forms\Components\TextInput::make('detail_name')
                            ->label('Name')
                            ->required()
                            ->maxLength(255),

                        Forms\Components\Textarea::make('detail_description')
                            ->label('Description')
                            ->required()
                            ->rows(3),

                        Forms\Components\TextInput::make('detail_point_gain')
                            ->label('Point Gain')
                            ->numeric()
                            ->required(),
                    ]),

                Forms\Components\Section::make('Event Details')
                    ->schema([
                        Forms\Components\TextInput::make('location')->required(),
                        Forms\Components\TextInput::make('contact')->required(),
                        Forms\Components\DatePicker::make('starts_at')->required(),
                        Forms\Components\DatePicker::make('ends_at')->required(),

                        Forms\Components\FileUpload::make('cover_url')
                            ->label('Cover Image')
                            ->disk('s3')
                            ->directory('events/covers')
                            ->image()
                            ->visibility('private')
                            ->required(),
                    ]),
            ]);
    }

    public static function table(Tables\Table $table): Tables\Table
    {
        return $table
            ->columns([
                Tables\Columns\TextColumn::make('detail.name')->label('Event Name'),
                Tables\Columns\TextColumn::make('location'),
                Tables\Columns\ImageColumn::make('cover_url')->label('Cover'),
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
            'index' => Pages\ListEvents::route('/'),
            'create' => Pages\CreateEvent::route('/create'),
            'edit' => Pages\EditEvent::route('/{record}/edit'),
        ];
    }
}
