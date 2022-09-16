<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Exercise extends Model
{
    use HasFactory;

    protected $fillable = [
        'user_id', 'name', 'weight', 'reps', 'sets', 'id'
    ];

    public function workout()
    {
        return $this->belongsToMany(Workout::class, 'exercise_workout');
    }
}
