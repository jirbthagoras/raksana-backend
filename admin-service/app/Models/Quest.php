<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Quest extends Model
{
    protected $fillable = [
            'detail_id',
            'code_id',
            'location',
            'latitude',
            'longitude',
            'max_contributors',
            'clue'
    ];

    public $timestamps = false;

    public function detail(): BelongsTo
    {
        return $this->belongsTo(Details::class, "detail_id");
    }

    public function code(): BelongsTo
    {
        return $this->belongsTo(Codes::class, "code_id");
    }
}
