<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasOne;

class Codes extends Model
{
    protected $fillable = [
        "id",
        "image_url"
    ];

    public $incrementing = false;   // UUID, not auto-increment
    protected $keyType = 'string';  // Important: id is string (uuid)

    public $timestamps = false;

    public function quest(): HasOne
    {
        return $this->hasOne(Quest::class, "code_id");
    }

    public function treasure(): HasOne
    {
        return $this->hasOne(Treasures::class);
    }

    public function event(): HasOne
    {
        return $this->hasOne(Event::class);
    }
}
