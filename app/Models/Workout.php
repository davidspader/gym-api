<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Workout extends Model
{
    use HasFactory;

    protected $fillable = [
        'id', 'name'
    ];

    public function exercises()
    {
        return $this->belongsToMany(Exercise::class, 'exercise_workout');
    }
}
