<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Event extends Model
{
    public function detail(): BelongsTo
    {
        return $this->belongsTo(Details::class, "detail_id");
    }

    public $timestamps = false;


    protected $fillable = [
    'detail_id',
    'code_id',
    'location',
    'contact',
    'starts_at',
    'ends_at',
    'cover_url',
    ];

    public function code(): BelongsTo
    {
        return $this->belongsTo(Codes::class, "code_id");
    }
}