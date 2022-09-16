<?php

namespace Tests\Feature\API\Exercise;

use App\Models\User;
use App\Models\Workout;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Illuminate\Testing\Fluent\AssertableJson;
use Tests\TestCase;

class WorkoutControllerTest extends TestCase
{
    use RefreshDatabase;

    private User $user;

    protected function setUp(): void
    {
        parent::setUp();
        $this->user = User::factory(1)->createOne();

        $response = $this->postJson('/api/auth/login', [
            'email' => $this->user->email,
            'password' => 'password',
        ]);
    }

    public function test_get_workout_endpoint()
    {
        $workout = Workout::factory(5)->create();

        $response = $this->getJson('api/workouts');

        $response->assertStatus(200);

        $response->assertJsonCount(5);

        $response->assertJson(function (AssertableJson $json) use ($workout) {

            $json->hasAll(['0.id', '0.user_id', '0.name', '0.exercises']);

            $json->whereAllType([
                '0.id' => 'integer',
                '0.user_id' => 'integer',
                '0.name' => 'string',
                '0.exercises' => 'array'
            ]);

            $workout = $workout->first();

            $json->whereAll([
                '0.id' => $workout->id,
                '0.user_id' => $workout->user_id,
                '0.name' => $workout->name,
                '0.exercises' => $workout->exercises
            ]);
        });
    }
}
