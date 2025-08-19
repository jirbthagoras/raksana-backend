<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Event extends Model
{
    public function detail(): BelongsTo
    {
        return $this->belongsTo(Details::class);
    }

    public function code(): BelongsTo
    {
        return $this->belongsTo(Codes::class);
    }
}