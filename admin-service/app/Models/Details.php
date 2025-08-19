<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasOne;

class Details extends Model
{
    public function quest(): HasOne
    {
        return $this->hasOne(Quest::class);
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
