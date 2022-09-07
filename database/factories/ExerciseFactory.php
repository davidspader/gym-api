<?php

namespace Database\Factories;

use Illuminate\Database\Eloquent\Factories\Factory;

/**
 * @extends \Illuminate\Database\Eloquent\Factories\Factory<\App\Models\Exercise>
 */
class ExerciseFactory extends Factory
{
    /**
     * Define the model's default state.
     *
     * @return array<string, mixed>
     */
    public function definition()
    {
        return [
            'user_id' => 1,
            'name' => $this->faker->sentence(3),
            'weight' => $this->faker->randomFloat(1, 12, 100),
            'reps' => $this->faker->numberBetween(4, 20),
            'sets' => $this->faker->numberBetween(2, 6)
        ];
    }
}
