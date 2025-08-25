<?php

namespace App\Filament\Resources;

use App\Filament\Resources\TreasureResource\Pages;
use App\Models\Treasure;
use App\Models\Treasures;
use Barryvdh\DomPDF\Facade\Pdf;
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
            ->headerActions([
                Tables\Actions\Action::make('export_recent_pdf')
                    ->label('Export Unclaimed Treasures')
                    ->icon('heroicon-o-document-text')
                    ->action(fn () => static::exportRecentPdf()),
            ])
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

    public static function exportRecentPdf()
        {
            $treasures = Treasures::with(['code'])
                ->where("claimed", "=", "false")
                ->get();

            if ($treasures->isEmpty()) {
                return redirect()->back()->with('danger', 'No quests created in the last 3 days.');
            }

            $pdf = Pdf::loadView('pdf.treasure', [
                'treasures' => $treasures,
            ])->setPaper('a4', 'portrait');
            $pdf->setOptions([
                'isHtml5ParserEnabled' => true,
                'isRemoteEnabled' => true,
            ]);

            return response()->streamDownload(
                fn () => print($pdf->output()),
                'unclaimed-treasures.pdf'
            );
    }
}
