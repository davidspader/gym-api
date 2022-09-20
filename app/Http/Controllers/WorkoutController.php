<?php

namespace App\Http\Controllers;

use App\Models\Workout;
use Illuminate\Http\Request;

class WorkoutController extends Controller
{
    private $user;

    public function __construct(private Workout $workout)
    {
        $this->user = auth('sanctum')->user();
    }

    public function index()
    {
        return $this->workout::with('exercises')->where('user_id', $this->user['id'])->get();
    }

    public function show($id)
    {
        $workout = $this->workout::with('exercises')->where('id', $id)->get();

        return response()->json($workout);
    }
}
