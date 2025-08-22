<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Treasures extends Model
{
    public function detail(): BelongsTo
    {
        return $this->belongsTo(Details::class);
    }

protected $fillable = [
        'name',
        'point_gain',
        'code_id',
        'claimed',
    ];

    public function code(): BelongsTo
    {
        return $this->belongsTo(Codes::class, 'code_id', 'id');
    }
}
