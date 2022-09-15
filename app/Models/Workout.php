<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Workout extends Model
{
    use HasFactory;

    public function workout()
    {
        return $this->belongsToMany(Exercise::class, 'exercises_workout');
    }
}
