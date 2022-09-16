<?php

namespace Database\Seeders;

use App\Models\ExerciseWorkout;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class ExerciseWorkoutSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        ExerciseWorkout::factory(10)->create();
    }
}
