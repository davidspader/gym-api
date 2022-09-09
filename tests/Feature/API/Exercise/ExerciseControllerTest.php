<?php

namespace Tests\Feature\API\Exercise;

use App\Models\Exercise;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Illuminate\Foundation\Testing\WithoutMiddleware;
use Illuminate\Testing\Fluent\AssertableJson;
use Tests\TestCase;

class ExerciseControllerTest extends TestCase
{
    use RefreshDatabase;

    private User $user;
    private $token;

    protected function setUp(): void
    {
        parent::setUp();
        $this->user = User::factory(1)->createOne();

        $response = $this->postJson('/api/auth/login', [
            'email' => $this->user->email,
            'password' => 'password',
        ]);

        $this->token = $response->json()['data']['token'];
    }

    public function test_get_exercises_endpoint()
    {
        $exercises = Exercise::factory(5)->create();

        $response = $this->getJson('api/exercises');

        $response->assertStatus(200);

        $response->assertJsonCount(5);

        $response->assertJson(function (AssertableJson $json) use ($exercises) {

            $json->hasAll(['0.id', '0.user_id', '0.name', '0.weight', '0.reps', '0.sets']);

            $json->whereAllType([
                '0.id' => 'integer',
                '0.user_id' => 'integer',
                '0.name' => 'string',
                '0.weight' => 'double|integer',
                '0.reps' => 'integer',
                '0.sets' => 'integer',
            ]);

            $exercise = $exercises->first();

            $json->whereAll([
                '0.id' => $exercise->id,
                '0.user_id' => $exercise->user_id,
                '0.name' => $exercise->name,
                '0.reps' => $exercise->reps,
                '0.sets' => $exercise->sets,
            ]);
        });
    }

    public function test_get_single_exercise_endpoint()
    {
        $exercise = Exercise::factory(1)->createOne();

        $response = $this->getJson('/api/exercises/1');

        $response->assertStatus(200);

        $response->assertJson(function (AssertableJson $json) use ($exercise) {

            $json->hasAll(['id', 'user_id', 'name', 'weight', 'reps', 'sets', 'created_at', 'updated_at']);

            $json->whereAllType([
                'id' => 'integer',
                'user_id' => 'integer',
                'name' => 'string',
                'weight' => 'double|integer',
                'reps' => 'integer',
                'sets' => 'integer',
            ]);

            $json->whereAll([
                'id' => $exercise->id,
                'user_id' => $exercise->user_id,
                'name' => $exercise->name,
                'reps' => $exercise->reps,
                'sets' => $exercise->sets,
            ]);
        });
    }

    public function test_get_wrong_user_id_in_single_exercise_endpoint()
    {
        User::factory(1)->createOne([
            "email" => 'test@test.com'
        ]);

        Exercise::factory(1)->createOne([
            'user_id' => 2
        ]);

        $response = $this->getJson('/api/exercises/1');

        $response->assertStatus(403);

        $response->assertJson(function (AssertableJson $json) {
            $json->hasAll('message')->etc();

            $json->whereAll([
                'message' => 'This action is unauthorized.'
            ]);
        });
    }

    public function test_post_exercise_endpoint()
    {
        $exercise = Exercise::factory(1)->makeOne()->toArray();

        $response = $this->postJson('/api/exercises', $exercise);

        $response->assertStatus(201);

        $response->assertJson(function (AssertableJson $json) use($exercise){

            $json->hasAll(['id', 'user_id', 'name', 'weight', 'reps', 'sets', 'created_at', 'updated_at']);

            $json->whereAll([
                'user_id' => $exercise['user_id'],
                'name' => $exercise['name'],
                'reps' => $exercise['reps'],
                'sets' => $exercise['sets'],
            ])->etc();
        });
    }

    public function test_get_wrong_user_id_in_post_exercise_endpoint()
    {
        User::factory(1)->createOne([
            "email" => 'test@test.com'
        ]);

        $exercise = Exercise::factory(1)->makeOne([
            'user_id' => 2
        ])->toArray();

        $response = $this->postJson('/api/exercises', $exercise);

        $response->assertStatus(403);

        $response->assertJson(function (AssertableJson $json) {
            $json->hasAll('message')->etc();

            $json->whereAll([
                'message' => 'This action is unauthorized.'
            ]);
        });
    }

    public function test_post_exercise_should_validate_when_try_create_a_invalid_exercise()
    {
        $response = $this->postJson('/api/exercises', []);

        $response->assertStatus(422);

        $response->assertJson(function (AssertableJson $json) {
            $json->hasAll(['message', 'errors']);

            $json->where('errors.user_id.0', 'The user id field is required.')
                ->where('errors.name.0', 'The name field is required.');
        });

        $exercise = Exercise::factory(1)->makeOne([
            'weight' => 'weight',
            'reps' => 'reps',
            'sets' => 'sets',
        ])->toArray();

        $response = $this->postJson('/api/exercises', $exercise);

        $response->assertStatus(422);

        $response->assertJson(function (AssertableJson $json) {
            $json->hasAll(['message', 'errors']);

            $json->where('errors.weight.0', 'The weight must be a number.')
                ->where('errors.reps.0', 'The reps must be an integer.')
                ->where('errors.sets.0', 'The sets must be an integer.');
        });

    }
}
