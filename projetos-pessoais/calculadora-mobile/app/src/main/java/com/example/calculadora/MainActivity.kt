package com.example.calculadora

import android.os.Bundle
import android.widget.Button
import android.widget.TextView
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.widget.AppCompatButton
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import com.example.calculadora.databinding.ActivityMainBinding
import org.mariuszgromada.math.mxparser.Expression

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        binding = ActivityMainBinding.inflate(layoutInflater)

        setContentView(binding.root)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        val calculo = binding.calculo
        val resultado = binding.resultado

        // Função para adicionar o número ao TextView
        fun adicionarNumero(numero: String) {
            calculo.text = "${calculo.text}$numero"
        }

        // Lista de pares (botão, valor)
        listOf(
            binding.numeroUm to "1",
            binding.numeroDois to "2",
            binding.numeroTres to "3",
            binding.numeroQuatro to "4",
            binding.numeroCinco to "5",
            binding.numeroSeis to "6",
            binding.numeroSete to "7",
            binding.numeroOito to "8",
            binding.numeroNove to "9",
            binding.numeroZero to "0",
            binding.parenteseAbrindo to "(",
            binding.parenteseFechando to ")",
            binding.divisao to "/",
            binding.multiplicacao to "*",
            binding.subtracao to "-",
            binding.soma to "+",
            binding.ponto to "."
        ).forEach { (botao, valor) ->
            botao.setOnClickListener {
                adicionarNumero(valor)
                val resultadoCalculado = Expression(calculo.text.toString()).calculate()

                resultado.text = resultadoCalculado.toString()
            }
        }

        binding.apagar.setOnClickListener {
            calculo.text = calculo.text.dropLast(1)
        }

        binding.ce.setOnClickListener { calculo.text = "" }

        binding.igualdade.setOnClickListener {
            val resultadoCalculado = Expression(calculo.text.toString()).calculate()

            resultado.text = resultadoCalculado.toString()
        }

    }
}