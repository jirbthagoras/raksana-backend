<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use League\CommonMark\Extension\CommonMark\Node\Inline\Code;

class Challenge extends Model
{
    public function detail(): BelongsTo
    {
        return $this->belongsTo(Details::class);
    }
}
