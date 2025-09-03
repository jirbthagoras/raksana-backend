<?php

namespace App\Filament\Resources;

use App\Filament\Resources\QuestResource\Pages;
use App\Models\Quest;
use Barryvdh\DomPDF\Facade\Pdf;
use Carbon\Carbon;
use Cheesegrits\FilamentGoogleMaps\Fields\Map;
use Filament\Forms;
use Filament\Tables;
use Filament\Forms\Form;
use Filament\Tables\Table;
use Filament\Resources\Resource;
use Filament\Tables\Columns\TextColumn;
use Illuminate\Support\Facades\Storage;

class QuestResource extends Resource
{
    protected static ?string $model = Quest::class;

    protected static ?string $navigationIcon = 'heroicon-o-flag';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                Forms\Components\Section::make('Quest Detail')
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

                // Quest Info fields
                Forms\Components\Section::make('Quest Info')
                    ->schema([
                        Forms\Components\TextInput::make('location')
                            ->label('Location Name')
                            ->required(),

                        Forms\Components\TextInput::make('max_contributors')
                            ->numeric()
                            ->required(),
                    ]),

                    Map::make('map_picker')
                        ->label('Quest Location (Pick on Map)')
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
                ->headerActions([
                Tables\Actions\Action::make('export_recent_pdf')
                    ->label('Export Last 3 Days')
                    ->icon('heroicon-o-document-text')
                    ->action(fn () => static::exportRecentPdf()),
            ])
            ->columns([
                Tables\Columns\TextColumn::make('detail.name')->label('Detail'),
                Tables\Columns\TextColumn::make('location'),
                Tables\Columns\TextColumn::make('max_contributors')->label('Max Contributors'),
                TextColumn::make('detail.point_gain')->label('Point Gain'),
                Tables\Columns\ImageColumn::make('code.image_url')->label('QR Code'),
            ])
            ->actions([
                Tables\Actions\EditAction::make(),
            ])
            ->bulkActions([
                Tables\Actions\DeleteBulkAction::make()
                    ->before(function ($records) {
                                foreach ($records as $record) {
                                    if ($record->code && $record->code->image_url) {
                                        $path = str_replace(Storage::disk('s3')->url(''), '', $record->code->image_url);

                                        if (Storage::disk('s3')->exists($path)) {
                                            Storage::disk('s3')->delete($path);                                        var_dump($path);
                                        }
                                    }
                                }
                            }),
            ]);
    }

    public static function getPages(): array
    {
        return [
            'index' => Pages\ListQuests::route('/'),
            'create' => Pages\CreateQuest::route('/create'),
            'edit' => Pages\EditQuest::route('/{record}/edit'),
        ];
    }

    public static function exportRecentPdf()
        {
            $quests = Quest::with(['detail', 'code'])
                ->whereHas('detail', function ($query) {
                    $query->where('created_at', '>=', Carbon::now()->subDays(3));
                })
                ->get();

            if ($quests->isEmpty()) {
                return redirect()->back()->with('danger', 'No quests created in the last 3 days.');
            }

            $pdf = Pdf::loadView('pdf.quest', [
                'quests' => $quests,
            ])->setPaper('a4', 'portrait');
            $pdf->setOptions([
                'isHtml5ParserEnabled' => true,
                'isRemoteEnabled' => true,
            ]);

            return response()->streamDownload(
                fn () => print($pdf->output()),
                'quests-last-3-days.pdf'
            );
        }
}
