<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasOne;

class Details extends Model
{
    protected $fillable = [
        "id",
        "name",
        "description",
        "point_gain"
    ];

    public function quest(): HasOne
    {
        return $this->hasOne(Quest::class, "detail_id");
    }

    public function treasure(): HasOne
    {
        return $this->hasOne(Treasures::class);
    }

    public function challenge(): HasOne
    {
        return $this->hasOne(Challenge::class);
    }

    public function event(): HasOne
    {
        return $this->hasOne(Event::class);
    }
}
